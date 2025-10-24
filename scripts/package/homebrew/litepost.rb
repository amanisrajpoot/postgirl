class Litepost < Formula
  desc "Lightweight Postman alternative"
  homepage "https://github.com/yourname/postgirl"
  version "0.1.0"
  
  if Hardware::CPU.arm?
    url "https://github.com/yourname/postgirl/releases/download/v0.1.0/postgirl-darwin-arm64.zip"
    sha256 "PUT_SHA256_HERE"
  else
    url "https://github.com/yourname/postgirl/releases/download/v0.1.0/postgirl-darwin-amd64.zip"
    sha256 "PUT_SHA256_HERE"
  end
  
  def install
    bin.install Dir["postgirl-darwin-*/postgirl"].first => "postgirl"
  end
  
  test do
    system "#{bin}/postgirl", "--version"
  end
end
