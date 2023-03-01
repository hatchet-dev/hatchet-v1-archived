import {
  Grid,
  H1,
  H2,
  HorizontalSpacer,
  P,
  GridCard,
} from "@hatchet-dev/hatchet-components";
import RunsList from "components/run/runslist";
import React from "react";

const HomeView: React.FunctionComponent = () => {
  return (
    <>
      <H1>Home</H1>
      <HorizontalSpacer spacepixels={16} />
      <H2>Popular Templates</H2>
      <HorizontalSpacer spacepixels={16} />
      <P>Most popular templates.</P>
      <HorizontalSpacer spacepixels={16} />
      <Grid>
        <GridCard>EKS</GridCard>
        <GridCard>RDS</GridCard>
        <GridCard>S3</GridCard>
      </Grid>
      <HorizontalSpacer spacepixels={24} />
      <H2>Recent Runs</H2>
    </>
  );
};

export default HomeView;
