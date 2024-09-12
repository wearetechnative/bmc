{ stdenv
, lib
, fetchFromGitHub
, bash
, pkgs
, pkgJsonifyAwsDotfiles
, makeWrapper
}:
stdenv.mkDerivation rec {

  pname = "bmc";
  version = "0.1.0";

  src = ./.;

  buildInputs = with pkgs; [
      awscli2
      aws-mfa
      bash
      granted
      jq
      dasel
      gum
      pkgJsonifyAwsDotfiles
    ];

  nativeBuildInputs = [ makeWrapper ];

  installPhase = ''
    mkdir -p $out/bin

    cp aws-profile-select.sh $out/bin/aws-profile-select.sh
    wrapProgram $out/bin/aws-profile-select.sh \
      --prefix PATH : ${lib.makeBinPath buildInputs }

    cp ecsconnect.sh $out/bin/ecsconnect.sh
    wrapProgram $out/bin/ecsconnect.sh \
      --prefix PATH : ${lib.makeBinPath buildInputs }

    cp ec2ls.sh $out/bin/ec2ls.sh
    wrapProgram $out/bin/ec2ls.sh \
      --prefix PATH : ${lib.makeBinPath buildInputs }

    cp profsel.sh $out/bin/profsel.sh
    wrapProgram $out/bin/profsel.sh \
      --prefix PATH : ${lib.makeBinPath buildInputs }
  '';
}
