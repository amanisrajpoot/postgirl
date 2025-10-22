class Litepost < Formula
  desc "Lightweight Postman alternative"
  homepage "https://github.com/yourname/litepost"
  version "0.1.0"
  
  if Hardware::CPU.arm?
    url "https://github.com/yourname/litepost/releases/download/v0.1.0/litepost-darwin-arm64.zip"
    sha256 "PUT_SHA256_HERE"
  else
    url "https://github.com/yourname/litepost/releases/download/v0.1.0/litepost-darwin-amd64.zip"
    sha256 "PUT_SHA256_HERE"
  end
  
  def install
    bin.install Dir["litepost-darwin-*/litepost"].first => "litepost"
  end
  
  test do
    system "#{bin}/litepost", "--version"
  end
end
