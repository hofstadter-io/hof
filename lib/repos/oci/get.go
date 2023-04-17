package oci

// func GetImage(src, format string) error {
// 	o := crane.GetOptions()
// 	imageMap := map[string]v1.Image{}
// 	indexMap := map[string]v1.ImageIndex{}
//
// 	ref, err := name.ParseReference(src, o.Name...)
// 	if err != nil {
// 		return nil
// 	}
//
// 	// Fetch the manifest using default credentials.
// 	rmt, err := remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
// 	if err != nil {
// 		return nil
// 	}
//
// 	// Prints the digest of registry.example.com/private/repo
// 	fmt.Println(rmt.Digest)
//
// 	// pull the entire index, not just a child image.
// 	if format == "oci" && rmt.MediaType.IsIndex() && o.Platform == nil {
// 		idx, err := rmt.ImageIndex()
// 		if err != nil {
// 			return err
// 		}
// 		indexMap[src] = idx
// 	}
//
// 	img, err := rmt.Image()
// 	if err != nil {
// 		return err
// 	}
// 	if cachePath != "" {
// 		img = cache.Image(img, cache.NewFilesystemCache(cachePath))
// 	}
// 	imageMap[src] = img
//
// 	switch format {
// 	case "tarball":
// 		if err := crane.MultiSave(imageMap, path); err != nil {
// 			return fmt.Errorf("saving tarball %s: %w", path, err)
// 		}
// 	case "legacy":
// 		if err := crane.MultiSaveLegacy(imageMap, path); err != nil {
// 			return fmt.Errorf("saving legacy tarball %s: %w", path, err)
// 		}
// 	case "oci":
// 		// Don't use crane.MultiSaveOCI so we can control annotations.
// 		p, err := layout.FromPath(path)
// 		if err != nil {
// 			p, err = layout.Write(path, empty.Index)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		for ref, img := range imageMap {
// 			opts := []layout.Option{}
// 			if annotateRef {
// 				parsed, err := name.ParseReference(ref, o.Name...)
// 				if err != nil {
// 					return err
// 				}
// 				opts = append(opts, layout.WithAnnotations(map[string]string{
// 					"org.opencontainers.image.ref.name": parsed.Name(),
// 				}))
// 			}
// 			if err = p.AppendImage(img, opts...); err != nil {
// 				return err
// 			}
// 		}
//
// 		for ref, idx := range indexMap {
// 			opts := []layout.Option{}
// 			if annotateRef {
// 				parsed, err := name.ParseReference(ref, o.Name...)
// 				if err != nil {
// 					return err
// 				}
// 				opts = append(opts, layout.WithAnnotations(map[string]string{
// 					"org.opencontainers.image.ref.name": parsed.Name(),
// 				}))
// 			}
// 			if err := p.AppendIndex(idx, opts...); err != nil {
// 				return err
// 			}
// 		}
// 	default:
// 		return fmt.Errorf("unexpected --format: %q (valid values are: tarball, legacy, and oci)", format)
// 	}
// 	return nil
// }
//
