# Maintainer: Andrey Kitsul <a.kitsul@zarplata.ru>
# Contributor: Andrey Kitsul <a.kitsul@zarplata.ru>

pkgname=zabbix-agent-extension-sentry
pkgver=${PKGVER:-autogenerated}
pkgrel=${PKGREL:-1}
_branch=${BRANCH:-master}
pkgdesc="Extension for zabbix-agentd for monitoring sentry"
arch=('any')
license=('GPL')
makedepends=('go')
depends=('zabbix-agent')
#install="install.sh"
source=("git+https://github.com/zarplata/$pkgname.git#branch=$_branch")
md5sums=(
    'SKIP'
    )

build() {
    cd "$srcdir/$pkgname"
    make 
}

package() {
	cd "$srcdir/$pkgname"
    ZBX_INC_DIR=/etc/zabbix/zabbix_agentd.conf.d/

    install -Dm 0755 .out/"${pkgname}" "${pkgdir}/usr/bin/${pkgname}"
    install -Dm 0644 "${pkgname}.conf" "${pkgdir}${ZBX_INC_DIR}${pkgname}.conf"
    
}
