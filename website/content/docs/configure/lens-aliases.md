---
title: Lens aliases
---

# Lens aliases

Camera manufactures often report various EXIF values for the same lens used on different camera bodies. This makes it difficult to count lens use totals across images and requires filtering under different labels even though the same lens was used to take the image.

For example, the Nikon AF-S Nikkor 70-200mm f/2.8G ED VR II is recorded under three values:

- AF-S Nikkor 70-200mm f/2.8G ED VR II
- 70.0-200.0 mm f/f2.8
- VR 70-200mm f/2.8G

rgallery features Lens aliases, allowing various values to be grouped under the same lens.

To enable lens aliases:

1. Create a file under the `config` directory called `config.yml`.
1. Populate it with the following YAML:

   ```yaml
   aliases:
     lenses:
       '70.0-200.0 mm f/f2.8': 'AF-S Nikkor 70-200mm f/2.8G ED VR II'
       'VR 70-200mm f/2.8G': 'AF-S Nikkor 70-200mm f/2.8G ED VR II'
       'AF-S Nikkor 70-200mm f/2.8G ED VR II': 'AF-S Nikkor 70-200mm f/2.8G ED VR II'
   ```

1. Restart the application.

All images with the lens value of `70.0-200.0 mm f/f2.8`, `VR 70-200mm f/2.8G`, `AF-S Nikkor 70-200mm f/2.8G ED VR II` will be totaled and can be filtered under the `AF-S Nikkor 70-200mm f/2.8G ED VR II` lens value (ex: `https://<replace-with-rgallery-url>/?lens=AF-S Nikkor 70-200mm f%2f2.8G ED VR II`).
