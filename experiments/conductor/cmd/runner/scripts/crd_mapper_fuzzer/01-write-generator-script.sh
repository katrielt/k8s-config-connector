#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

cd $(dirname "$0")
SCRIPT_DIR=`pwd`

PROMPT=${SCRIPT_DIR}/01-generate-script.prompt

if [[ -z "${WORKDIR}" ]]; then
  echo "WORKDIR is required"
  exit 1
fi

if [[ -z "${BRANCH_NAME}" ]]; then
  echo "BRANCH_NAME is required"
  exit 1
fi

if [[ -z "${LOG_DIR}" ]]; then
  echo "LOG_DIR is required"
  exit 1
fi


mkdir -p ${LOG_DIR}

cd ${WORKDIR}

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd ${REPO_ROOT}

git co master
git co ${BRANCH_NAME} || git co -b ${BRANCH_NAME}

mkdir -p apis/${SERVICE}/${CRD_VERSION}
cd apis/${SERVICE}/${CRD_VERSION}

if [[ ! -f generate.sh ]]; then

cat > generate.sh <<EOF
#!/bin/bash
# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="\$(git rev-parse --show-toplevel)"
cd \${REPO_ROOT}/dev/tools/controllerbuilder
EOF

fi


cat >> generate.sh <<EOF

go run . generate-types \
    --service ${PROTO_PACKAGE} \
    --api-version ${CRD_GROUP}/${CRD_VERSION} \
    --resource ${CRD_KIND}:${PROTO_RESOURCE}

go run . generate-mapper \
    --service ${PROTO_PACKAGE} \
    --api-version ${CRD_GROUP}/${CRD_VERSION}


cd \${REPO_ROOT}
dev/tasks/generate-crds

go run -mod=readonly golang.org/x/tools/cmd/goimports@latest -w  pkg/controller/direct/${SERVICE}/

EOF

chmod +x generate.sh

git status
git add .
git commit -m "${CRD_KIND}: generated scripts"

cd ${REPO_ROOT}

./dev/tools/controllerbuilder/generate-proto.sh
./apis/${SERVICE}/${CRD_VERSION}/generate.sh


git status
git add .
git commit -m "autogen: apis/${SERVICE}/${CRD_VERSION}/generate.sh"

echo "Done"