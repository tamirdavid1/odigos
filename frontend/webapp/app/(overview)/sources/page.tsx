'use client';

import React from 'react';
import { TABLE_MAX_HEIGHT, TABLE_MAX_WIDTH } from '@/utils';
import { useMetrics, useSourceCRUD, useWorkloadUtils } from '@/hooks';
import { MultiSourceControl, SourceTable } from '@odigos/ui-kit/containers';

export default function Page() {
  const { metrics } = useMetrics();
  const { restartWorkloads } = useWorkloadUtils();
  const { sources, persistSources } = useSourceCRUD();

  return (
    <>
      <SourceTable metrics={metrics} maxHeight={TABLE_MAX_HEIGHT} maxWidth={TABLE_MAX_WIDTH} />
      <MultiSourceControl totalSourceCount={sources.length} uninstrumentSources={(payload) => persistSources(payload, {})} restartWorkloads={restartWorkloads} />
    </>
  );
}
