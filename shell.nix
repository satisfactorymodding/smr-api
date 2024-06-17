{ pkgs, unstable }:

pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    libwebp
    libpng
    unstable.go_1_22
    protobuf
    protoc-gen-go-grpc
    protoc-gen-go
    minio-client
    unstable.golangci-lint
    unstable.delve
  ];
}
