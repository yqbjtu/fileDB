export APP_NAME=fileDB
export BUILD_VERSION=1.0.0
export BUILD_IMAGE=fileDB/${APP_NAME}


echo "docker build -t ${BUILD_IMAGE}:${BUILD_VERSION} ."
docker build -t "${BUILD_IMAGE}:${BUILD_VERSION}" .

#下面可以加上docker push的命令。 把镜像上传到能访问的docker image registry，因为本人都是本地（本地k8s采用的docker save，然后docker load镜像的模式）验证测试
