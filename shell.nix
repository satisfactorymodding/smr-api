{ pkgs, unstable }:

pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    libwebp
    libpng
    unstable.go_1_21
    protobuf
    protoc-gen-go-grpc
    minio-client
    unstable.golangci-lint
    unstable.delve
  ];
}
