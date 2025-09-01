local rgallery = import 'rgallery/rgallery.libsonnet';
local k = import 'ksonnet-util/kausal.libsonnet';

rgallery {
  _config+:: {
    namespace: 'rgallery',
  },
  namespace: k.core.v1.namespace.new($._config.namespace),
  rgallery: rgallery.new(name='demo', image='robbymilo/rgallery:latest', dataPath='/mnt/rgallery/data', mediaPath='/mnt/rgallery/media', cachePath='/mnt/rgallery/cache') {
  },
}
