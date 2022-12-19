import BackText from "components/backtext";
import Breadcrumbs from "components/breadcrumbs";
import {
  FlexRowRight,
  H1,
  H2,
  H3,
  H4,
  HorizontalSpacer,
  P,
  Span,
} from "components/globals";
import Paginator from "components/paginator";
import SectionCard from "components/sectioncard";
import Table from "components/table";
import TabList from "components/tablist";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";

type Props = {
  back: () => void;
};

const ExpandedRun: React.FC<Props> = ({ back }) => {
  return (
    <>
      <BackText text="All Runs" back={back} />
      <HorizontalSpacer spacepixels={12} />
      <SectionCard>
        <H4>Overview</H4>
      </SectionCard>
      <HorizontalSpacer spacepixels={12} />
      <SectionCard>
        <H4>Configuration</H4>
      </SectionCard>
      <HorizontalSpacer spacepixels={12} />
      <SectionCard>
        <H4>Logs</H4>
      </SectionCard>
    </>
  );
};

export default ExpandedRun;
