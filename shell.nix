{ pkgs ? import <nixpkgs> {} }:

let
  unstable = import (fetchTarball https://nixos.org/channels/nixos-unstable/nixexprs.tar.xz) { };
in
pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    libwebp
    libpng
    unstable.go_1_21
    protobuf
    protoc-gen-go-grpc
    minio-client
    unstable.golangci-lint
  ];
}
