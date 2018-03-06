{ mkDerivation, base, stdenv }:
mkDerivation {
  pname = "vrf-hs";
  version = "0.1.0.0";
  src = ./.;
  libraryHaskellDepends = [ base ];
  librarySystemDepends = [ ];
  homepage = "https://github.com/sakshamsharma/vrf-hs#readme";
  license = stdenv.lib.licenses.bsd3;
}
