{ pkgs ? import <nixpkgs> {} }:

with pkgs;

mkShell {
  buildInputs = [
        #go
        #gopls
        runit
  ];
}
