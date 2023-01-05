import Breadcrumbs from "components/breadcrumbs";
import {
  FlexRowRight,
  Grid,
  H1,
  H2,
  HorizontalSpacer,
  P,
  Span,
} from "components/globals";
import { GridCard } from "components/gridcard";
import Example from "components/heirarchygraph";
import Paginator from "components/paginator";
import RunsList from "components/runslist";
import Table from "components/table";
import TabList from "components/tablist";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";

const TabOptions = ["Runs", "Resource Explorer", "Configuration", "Settings"];

const HomeView: React.FunctionComponent = () => {
  const [selectedTab, setSelectedTab] = useState(TabOptions[0]);

  let history = useHistory();

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
      <RunsList
        runs={[
          {
            status: "deployed",
            date: "7:09 AM on June 23rd, 2022",
          },
        ]}
      />
    </>
  );
};

export default HomeView;
