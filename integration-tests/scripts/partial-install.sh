#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail
set -x

echo "[INFO]Configuring deployment"

# Variables for RHTAP Sample Backstage Templates
DEVELOPER_HUB__CATALOG__URL="https://github.com/redhat-appstudio/tssc-sample-templates/blob/main/all.yaml"
# Variables for GitHub integration
GITHUB__APP__ID=$(cat /usr/local/rhtap-cli-install/rhdh-github-app-id)
GITHUB__APP__CLIENT__ID=$(cat /usr/local/rhtap-cli-install/rhdh-github-client-id)
GITHUB__APP__CLIENT__SECRET=$(cat /usr/local/rhtap-cli-install/rhdh-github-client-secret)
GITHUB__APP__PRIVATE_KEY=$(base64 -d < /usr/local/rhtap-cli-install/rhdh-github-private-key | sed 's/^/        /')
GITOPS__GIT_TOKEN=$(cat /usr/local/rhtap-cli-install/github_token)
GITHUB__APP__WEBHOOK__SECRET=$(cat /usr/local/rhtap-cli-install/rhdh-github-webhook-secret)
# Variables for Gitlab integration
GITLAB__TOKEN=$(cat /usr/local/rhtap-cli-install/gitlab_token)
# Variables for Jenkins integration
JENKINS_API_TOKEN=$(cat /usr/local/rhtap-cli-install/jenkins-api-token)
JENKINS_URL=$(cat /usr/local/rhtap-cli-install/jenkins-url)
JENKINS_USERNAME=$(cat /usr/local/rhtap-cli-install/jenkins-username)
## Variables for quay.io integration
QUAY__DOCKERCONFIGJSON=$(cat /usr/local/rhtap-cli-install/quay-dockerconfig-json)
QUAY__API_TOKEN=$(cat /usr/local/rhtap-cli-install/quay-api-token)
## Variables for ACS integration
ACS__CENTRAL_ENDPOINT=$(cat /usr/local/rhtap-cli-install/acs-central-endpoint)
ACS__API_TOKEN=$(cat /usr/local/rhtap-cli-install/acs-api-token)

readonly tpl_file="installer/charts/values.yaml.tpl"
readonly config_file="installer/config.yaml"

ci_enabled() {
  echo "[INFO] Turn ci to true, this is required when you perform rhtap-e2e automation test against RHTAP"
  sed -i'' -e 's/ci: false/ci: true/g' "$tpl_file"
}

update_dh_catalog_url() {
  echo "[INFO] Update dh catalog url with $DEVELOPER_HUB__CATALOG__URL"
  yq -i ".rhtapCLI.features.redHatDeveloperHub.properties.catalogURL = strenv(DEVELOPER_HUB__CATALOG__URL)" "${config_file}"
}

github_integration() {
  echo "[INFO] Config Github App integration in RHTAP"

  cat <<EOF >>"$tpl_file"
integrations:
  github:
    id: "${GITHUB__APP__ID}"
    clientId: "${GITHUB__APP__CLIENT__ID}"
    clientSecret: "${GITHUB__APP__CLIENT__SECRET}"
    host: "github.com"
    publicKey: |-
$(echo "${GITHUB__APP__PRIVATE_KEY}" | sed 's/^/      /')
    token: "${GITOPS__GIT_TOKEN}"
    webhookSecret: "${GITHUB__APP__WEBHOOK__SECRET}"
EOF
}

jenkins_integration() {
  echo "[INFO] Integrates an exising Jenkins server into RHTAP"
  ./bin/rhtap-cli integration --kube-config "$KUBECONFIG" jenkins --token="$JENKINS_API_TOKEN" --url="$JENKINS_URL" --username="$JENKINS_USERNAME" --force
}

gitlab_integration() {
  echo "[INFO] Configure an external Gitlab integration into RHTAP"
  ./bin/rhtap-cli integration --kube-config "$KUBECONFIG" gitlab --token "${GITLAB__TOKEN}"
}

quay_integration() {
  echo "[INFO] Configure quay.io integration into RHTAP"
  ./bin/rhtap-cli integration --kube-config "$KUBECONFIG" quay --url="https://quay.io" --dockerconfigjson="${QUAY__DOCKERCONFIGJSON}" --token="${QUAY__API_TOKEN}"
}

acs_integration() {
  echo "[INFO] Configure an external ACS integration into RHTAP"
  ./bin/rhtap-cli integration --kube-config "$KUBECONFIG" acs --endpoint="${ACS__CENTRAL_ENDPOINT}" --token="${ACS__API_TOKEN}"
}

install_rhtap() {
  echo "[INFO] Start installing RHTAP"
  github_integration
  echo "[INFO] Building binary"
  make build

  echo "[INFO] Installing RHTAP"
  jenkins_integration
  gitlab_integration
  quay_integration
  acs_integration
  ./bin/rhtap-cli deploy --timeout 30m --embedded false --config "$config_file" --values-template "$tpl_file" --kube-config "$KUBECONFIG"

  homepage_url=https://$(kubectl -n rhtap get route backstage-developer-hub -o  'jsonpath={.spec.host}')
  callback_url=https://$(kubectl -n rhtap get route backstage-developer-hub -o  'jsonpath={.spec.host}')/api/auth/github/handler/frame
  webhook_url=https://$(kubectl -n openshift-pipelines get route pipelines-as-code-controller -o 'jsonpath={.spec.host}')

  echo "[INFO]homepage_url=$homepage_url"
  echo "[INFO]callback_url=$callback_url"
  echo "[INFO]webhook_url=$webhook_url"
}

ci_enabled
update_dh_catalog_url
install_rhtap