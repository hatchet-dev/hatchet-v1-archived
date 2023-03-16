#!/bin/bash

version=$1

create_hatchet_formula() {
    name=hatchet_${version}_Darwin_x86_64.zip
    curl -L https://github.com/hatchet-dev/hatchet/releases/download/${version}/${name} --output $name

    sha=$(cat ${name} | openssl dgst -sha256 -binary | xxd -p -c 256)

    cat >hatchet.rb <<EOL
class Hatchet < Formula
    depends_on "hatchet-server"
    depends_on "hatchet-admin"

    homepage "https://hatchet.run"
    version "${version}"

    url "https://github.com/hatchet-dev/hatchet/releases/download/${version}/${name}" 
    sha256 "${sha}"
          
    on_macos do
      def install
        bin.install "hatchet"
      end
    end
  end
EOL

    rm $name
}

create_hatchet_server_formula() {
    name=hatchet-server_${version}_Darwin_x86_64.zip
    curl -L https://github.com/hatchet-dev/hatchet/releases/download/${version}/${name} --output $name

    sha=$(cat ${name} | openssl dgst -sha256 -binary | xxd -p -c 256)

    cat >hatchet-server.rb <<EOL
class HatchetServer < Formula
    homepage "https://hatchet.run"
    version "${version}"

    url "https://github.com/hatchet-dev/hatchet/releases/download/${version}/${name}" 
    sha256 "${sha}"
          
    on_macos do
      def install
        bin.install "hatchet-server"
      end
    end
  end
EOL

    rm $name
}

create_hatchet_admin_formula() {
    name=hatchet-admin_${version}_Darwin_x86_64.zip
    curl -L https://github.com/hatchet-dev/hatchet/releases/download/${version}/${name} --output $name

    sha=$(cat ${name} | openssl dgst -sha256 -binary | xxd -p -c 256)

    cat >hatchet-admin.rb <<EOL
class HatchetAdmin < Formula
    homepage "https://hatchet.run"
    version "${version}"

    url "https://github.com/hatchet-dev/hatchet/releases/download/${version}/${name}" 
    sha256 "${sha}"
          
    on_macos do
      def install
        bin.install "hatchet-admin"
      end
    end
  end
EOL

    rm $name
}

create_hatchet_formula
create_hatchet_server_formula
create_hatchet_admin_formula