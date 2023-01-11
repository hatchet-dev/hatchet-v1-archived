import BackText from "components/backtext";
import { H4, HorizontalSpacer } from "components/globals";
import SectionCard from "components/sectioncard";
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
