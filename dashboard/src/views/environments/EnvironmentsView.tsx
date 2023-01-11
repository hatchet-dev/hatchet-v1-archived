import { Grid, H1, HorizontalSpacer, P } from "components/globals";
import { GridCard } from "components/gridcard";
import React from "react";
import { useHistory } from "react-router-dom";

const EnvironmentsView: React.FunctionComponent = () => {
  let history = useHistory();

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
