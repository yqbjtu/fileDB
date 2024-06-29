export APP_NAME=fileDB
export BUILD_VERSION=1.0.0
export BUILD_IMAGE=gostudy/${APP_NAME}



echo "docker build -t ${BUILD_IMAGE}:${BUILD_VERSION} ."
docker build -t "${BUILD_IMAGE}:${BUILD_VERSION}" .
