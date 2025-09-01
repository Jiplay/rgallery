{
  new(name, image, dataPath, mediaPath, cachePath):: {
    local k = (import 'ksonnet-util/kausal.libsonnet'),
    local this = self,
    local statefulset = k.apps.v1.statefulSet,
    local service = k.core.v1.service,
    local configmap = k.core.v1.configMap,
    local container = k.core.v1.container,
    local containerPort = k.core.v1.containerPort,
    local persistentVolume = k.core.v1.persistentVolume,
    local persistentVolumeClaim = k.core.v1.persistentVolumeClaim,
    local volume = k.core.v1.volume,
    local volumeMount = k.core.v1.volumeMount,
    local envVar = k.core.v1.envVar,
    local utils = import 'utils/utils.libsonnet',

    _config+: {
      rgallery_config_template: importstr 'files/rgallery.yml',
    },

    rgallery_container::
      container.new(name, image)
      + container.withPorts(
        [
          containerPort.new(name, 3000),
          containerPort.new('metrics', 3001),
        ]
      )
      + container.withEnvMixin([
        envVar.new('TZ', 'Europe/Stockholm'),
      ])
      + container.withCommandMixin(['/usr/bin/rgallery', '--include-originals=true'])
      + utils.buildHealthcheck(port=3000, periodSeconds=10, initialDelaySeconds=1, timeoutSeconds=10)
      + container.withVolumeMounts([
        volumeMount.new(name + '-cache', '/cache'),
        volumeMount.new(name + '-media', '/media'),
        volumeMount.new(name + '-data', '/data'),
      ])
      + container.securityContext.withRunAsUser(1000)
      + container.securityContext.withRunAsGroup(1000),

    rgallery_config_template:
      configmap.new('rgallery-config')
      + configmap.withData({
        'config.yml': std.format(this._config.rgallery_config_template, this._config),
      }),

    rgallery_statefulset:
      statefulset.new(name, 1, [this.rgallery_container])
      + statefulset.mixin.spec.template.spec.withVolumes([
        {
          name: name + '-media',
          persistentVolumeClaim: {
            claimName: name + '-media',
          },
        },
        {
          name: name + '-cache',
          persistentVolumeClaim: {
            claimName: name + '-cache',
          },
        },
        {
          name: name + '-data',
          persistentVolumeClaim: {
            claimName: name + '-data',
          },
        },
      ])
      + statefulset.spec.withServiceName(this.rgallery_statefulset.metadata.name)
      + statefulset.mixin.spec.template.metadata.withAnnotationsMixin({
        configVersion: std.md5(std.format(this._config.rgallery_config_template, this._config)),
      })
      + k.util.configVolumeMount('rgallery-config', '/config'),

    rgallery_service:
      k.util.serviceFor(this.rgallery_statefulset)
      + service.mixin.spec.withType('ClusterIP'),

    rgallery_data_volume:
      persistentVolume.new(name + '-data')
      + persistentVolume.spec.withAccessModes(['ReadWriteOnce'])
      + persistentVolume.spec.withCapacity({ storage: '20Gi' })
      + persistentVolume.spec.withStorageClassName('')
      + persistentVolume.spec.hostPath.withPath(dataPath),

    rgallery_data_claim:
      persistentVolumeClaim.new(name + '-data')
      + persistentVolumeClaim.spec.withAccessModes(['ReadWriteOnce'])
      + persistentVolumeClaim.spec.resources.withRequests({ storage: '20Gi' })
      + persistentVolumeClaim.spec.withStorageClassName('')
      + persistentVolumeClaim.spec.withVolumeName(name + '-data'),

    rgallery_media_volume:
      persistentVolume.new(name + '-media')
      + persistentVolume.spec.withAccessModes(['ReadWriteOnce'])
      + persistentVolume.spec.withCapacity({ storage: '20Gi' })
      + persistentVolume.spec.withStorageClassName('')
      + persistentVolume.spec.hostPath.withPath(mediaPath),

    rgallery_media_claim:
      persistentVolumeClaim.new(name + '-media')
      + persistentVolumeClaim.spec.withAccessModes(['ReadWriteOnce'])
      + persistentVolumeClaim.spec.resources.withRequests({ storage: '20Gi' })
      + persistentVolumeClaim.spec.withStorageClassName('')
      + persistentVolumeClaim.spec.withVolumeName(name + '-media'),

    rgallery_cache_volume:
      persistentVolume.new(name + '-cache')
      + persistentVolume.spec.withAccessModes(['ReadWriteOnce'])
      + persistentVolume.spec.withCapacity({ storage: '20Gi' })
      + persistentVolume.spec.withStorageClassName('')
      + persistentVolume.spec.hostPath.withPath(cachePath),

    rgallery_cache_claim:
      persistentVolumeClaim.new(name + '-cache')
      + persistentVolumeClaim.spec.withAccessModes(['ReadWriteOnce'])
      + persistentVolumeClaim.spec.resources.withRequests({ storage: '20Gi' })
      + persistentVolumeClaim.spec.withStorageClassName('')
      + persistentVolumeClaim.spec.withVolumeName(name + '-cache'),

  },
}
