import {
  H4,
  HorizontalSpacer,
  BackText,
  SectionCard,
} from "@hatchet-dev/hatchet-components";
import React from "react";

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
