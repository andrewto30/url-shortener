{
  pkgs,
  lib,
  config,
  ...
}:
{
  # https://devenv.sh/languages/
  languages.go.enable = true;

  # https://devenv.sh/services/
  services.redis.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}

