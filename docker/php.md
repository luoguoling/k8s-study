```dock
FROM php:7.2-fpm
LABEL maintainer="rolin"

#Download PHP extensions
#ADD https://raw.githubusercontent.com/mlocati/docker-php-extension-installer/master/install-php-extensions /usr/local/bin/
#RUN chmod uga+x /usr/local/bin/install-php-extensions && sync
RUN set -eux; \
    mv /etc/apt/sources.list /etc/apt/sources.list.bak \
    && echo "deb http://mirrors.aliyun.com/debian/ buster main non-free contrib \n \
      deb-src http://mirrors.aliyun.com/debian/ buster main non-free contrib \n \
      deb http://mirrors.aliyun.com/debian-security buster/updates main \n \
      deb-src http://mirrors.aliyun.com/debian-security buster/updates main \n \
      deb http://mirrors.aliyun.com/debian/ buster-updates main non-free contrib \n \
      deb-src http://mirrors.aliyun.com/debian/ buster-updates main non-free contrib \n \
      deb http://mirrors.aliyun.com/debian/ buster-backports main non-free contrib \n \
      deb-src http://mirrors.aliyun.com/debian/ buster-backports main non-free contrib" > /etc/apt/sources.list
COPY --from=mlocati/php-extension-installer /usr/bin/install-php-extensions /usr/bin/

RUN DEBIAN_FRONTEND=noninteractive apt-get update -q \
    && DEBIAN_FRONTEND=noninteractive apt-get install -qq -y \
      curl \
      git  \
      zip unzip \
    && install-php-extensions \
      bcmath \
      bz2 \
      calendar \
      exif \
      gd \
      intl \
      ldap \
      memcached \
      mysqli \
      opcache \
      pdo_mysql \
      pdo_pgsql \
      pgsql \
      redis \
      soap \
      xsl \
      zip \
      sockets \
      swoole \
      yaf \
      memcached \
      mongodb \
      mcrypt \
      iconv \
      mbstring
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && ln -s $(composer config --global home) /root/composer
ENV PATH=$PATH:/root/composer/vendor/bin COMPOSER_ALLOW_SUPERUSER=1
```

### 参考2 

> https://github.com/voocel/docker-lnmp/edit/master/php

```yaml
#此方法能够正常处理
FROM php:7.2-fpm

RUN set -eux; \
    mv /etc/apt/sources.list /etc/apt/sources.list.bak \
    && echo "deb http://mirrors.aliyun.com/debian/ buster main non-free contrib \n \
      deb-src http://mirrors.aliyun.com/debian/ buster main non-free contrib \n \
      deb http://mirrors.aliyun.com/debian-security buster/updates main \n \
      deb-src http://mirrors.aliyun.com/debian-security buster/updates main \n \
      deb http://mirrors.aliyun.com/debian/ buster-updates main non-free contrib \n \
      deb-src http://mirrors.aliyun.com/debian/ buster-updates main non-free contrib \n \
      deb http://mirrors.aliyun.com/debian/ buster-backports main non-free contrib \n \
      deb-src http://mirrors.aliyun.com/debian/ buster-backports main non-free contrib" > /etc/apt/sources.list
#如果是7.4需要重新安装libzip      
RUN set -eux; \
    apt-get update && apt-get install -y \
       git vim wget unzip bzip2  libbz2-dev libjpeg-dev libpng-dev curl libcurl4-openssl-dev libonig-dev \
       libmagickwand-dev libmcrypt-dev libonig-dev libxml2-dev libfreetype6-dev libjpeg62-turbo-dev zlib1g-dev \
    && docker-php-ext-install -j$(nproc) gd \
    && docker-php-ext-install zip pdo_mysql opcache mysqli mbstring bz2 soap bcmath calendar exif gettext sockets pcntl \
    && pecl install imagick-3.4.4 mcrypt-1.0.4 redis-5.3.4 xdebug-3.0.4 swoole-4.6.7 protobuf grpc \
    && docker-php-ext-enable imagick mcrypt redis xdebug swoole protobuf grpc \
    && rm -r /var/lib/apt/lists/* \
    && usermod -u 1000 www-data \
    && groupmod -g 1000 www-data

ENV COMPOSER_HOME /root/composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer
    
    
    

```

