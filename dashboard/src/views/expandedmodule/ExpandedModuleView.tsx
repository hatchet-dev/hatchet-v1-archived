import {
  H1,
  HorizontalSpacer,
  P,
  Span,
  Breadcrumbs,
  HeirarchyGraph,
  TabList,
  StandardButton,
  FlexRowRight,
  FlexCol,
} from "@hatchet-dev/hatchet-components";
import { useMutation } from "@tanstack/react-query";
import { useAtom } from "jotai";
import React, { useState, useEffect } from "react";
import { useHistory, useParams } from "react-router-dom";
import api from "shared/api";
import { currTeamAtom } from "shared/atoms/atoms";
import RunsList from "../../components/runslist";
import ModuleRunsList from "./components/modulerunslist";
import ExpandedModuleMonitors from "./components/monitors";

const TabOptions = [
  "Runs",
  "Resource Explorer",
  "Policies",
  "Monitors",
  "Settings",
];

const ExpandedModuleView: React.FunctionComponent = () => {
  const history = useHistory();
  const [selectedTab, setSelectedTab] = useState(TabOptions[0]);
  const [currTeam] = useAtom(currTeamAtom);
  const [err, setErr] = useState("");
  const params: any = useParams();

  useEffect(() => {
    if (!params?.module) {
      history.push(`/team/${currTeam?.id}/modules`);
    }
  }, [params]);

  const mutation = useMutation({
    mutationKey: ["create_module_run", currTeam?.id, params?.module],
    mutationFn: () => {
      return api.createModuleRun(currTeam?.id, params?.module);
    },
    onSuccess: (data) => {
      setErr("");
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const renderTabContents = () => {
    switch (selectedTab) {
      case "Runs":
        return (
          <FlexCol height="100%">
            <FlexRowRight>
              <StandardButton
                label="New Run"
                on_click={() => {
                  mutation.mutate();
                }}
              />
            </FlexRowRight>

            <ModuleRunsList team_id={currTeam.id} module_id={params?.module} />
          </FlexCol>
        );
      case "Resource Explorer":
        return <HeirarchyGraph width={100} height={100} />;
      case "Monitors":
        return <ExpandedModuleMonitors />;
      default:
        return <Span>Settings</Span>;
    }
  };

  return (
    <>
      <Breadcrumbs
        breadcrumbs={[
          {
            label: "Modules",
            link: "/modules",
          },
          {
            label: "Staging: team-1-gke",
            link: "",
          },
        ]}
      />
      <HorizontalSpacer spacepixels={12} />
      <H1>Staging: team-1-gke</H1>
      <HorizontalSpacer spacepixels={20} />
      <P>
        This page contains information about the team-1-gke workspace in the
        Staging environment.
      </P>
      <HorizontalSpacer spacepixels={20} />
      <TabList tabs={TabOptions} selectTab={setSelectedTab} />
      {renderTabContents()}
    </>
  );
};

export default ExpandedModuleView;
