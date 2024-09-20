{ stdenv
, lib
, fetchFromGitHub
, bash
, pkgs
, makeWrapper
}:
stdenv.mkDerivation rec {

  pname = "bmc";
  version = "0.1.0";

  src = ./.;

  buildInputs = with pkgs; [
      bash
    ];

  nativeBuildInputs = [ makeWrapper ];


  installPhase = ''
    mkdir -p $out/bin

    cp VERSION-bmc bmc *.sh $out/bin/
    runHook postInstall

#    wrapProgram $out/bin/aws-profile-select.sh --prefix PATH : ${lib.makeBinPath buildInputs }
#    wrapProgram $out/bin/profsel.sh --prefix PATH : ${lib.makeBinPath buildInputs }
#    wrapProgram $out/bin/ecsconnect.sh --prefix PATH : ${lib.makeBinPath buildInputs }
#    wrapProgram $out/bin/ec2ls.sh --prefix PATH : ${lib.makeBinPath buildInputs }
#    wrapProgram $out/bin/ec2connect.sh --prefix PATH : ${lib.makeBinPath buildInputs }

  '';
}
