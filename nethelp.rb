class Nethelp < Formula
  desc "Find out why you can't reach Sauce Lab's services. (Real Device Cloud, Virtual Cloud, Sauce Connect and more)"
  homepage "https://github.com/mdsauce/nethelp"
  url "https://github.com/mdsauce/nethelp/archive/v1.1.tar.gz"
  sha256 "669e3f8f6f1c4b80e4fe045d1202c7236a576643ea8413d6934fd9fee8e90e6f"
  depends_on "go" => :build
  version "1.1"

  def install
    system "go", "build", "-o", bin/"nethelp", "."
  end

  test do
    system "#{bin}/nethelp", --help, ">", "output.txt"
    assert_predicate ./"output.txt", :exist?
  end
end
