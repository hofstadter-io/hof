import { Metadata, ResolvedMetadata } from 'next';
import { deepmerge } from 'deepmerge-ts';

type CreateMetadataOpts = {
  title: string;
  description?: string;
  path: string;
};

export const createMetadata = (
  opts: CreateMetadataOpts,
  parentMetadata?: ResolvedMetadata,
  extras?: Metadata
): Metadata => {
  const title = opts.title;
  const description =
    opts.description ?? parentMetadata?.description ?? undefined;

  const base: Metadata = {
    // For some reason Next.js is returning a very weird object for `metadataBase`
    // from the parentMetadata, so we need to always set it here even if it is
    // already set in app/layout.tsx
    metadataBase: new URL('https://rafaelalmeidatk.com'),
    title,
    description,
    alternates: {
      canonical: opts.path,
    },
    openGraph: {
      title,
      description,
      url: opts.path,
    },
    twitter: {
      title,
      description,
    },
  };

  return deepmerge(parentMetadata, base, extras ?? {});
};
