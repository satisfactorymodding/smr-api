{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    libwebp
    go
    protobuf
    protoc-gen-go-grpc
    minio-client
  ];
}
