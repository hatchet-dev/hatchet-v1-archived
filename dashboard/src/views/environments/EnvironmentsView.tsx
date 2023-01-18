import {
  Grid,
  H1,
  HorizontalSpacer,
  P,
  GridCard,
} from "@hatchet-dev/hatchet-components";
import React from "react";

const EnvironmentsView: React.FunctionComponent = () => {
  return (
    <>
      <H1>Environments</H1>
      <HorizontalSpacer spacepixels={12} />
      <P>View and manage your environments.</P>
      <HorizontalSpacer spacepixels={20} />
      <Grid>
        <GridCard>Dev</GridCard>
        <GridCard>Staging</GridCard>
        <GridCard>Production</GridCard>
      </Grid>
    </>
  );
};

export default EnvironmentsView;
