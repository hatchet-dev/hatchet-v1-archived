import {
  FlexRowRight,
  H4,
  HorizontalSpacer,
  StandardButton,
  SectionArea,
  Table,
} from "hatchet-components";
import ResultsList from "components/monitor/resultslist";
import React from "react";

type Props = {
  team_id: string;
  module_id: string;
};

const ExpandedModuleMonitors: React.FC<Props> = ({ team_id, module_id }) => {
  return <ResultsList team_id={team_id} module_id={module_id} />;
};

export default ExpandedModuleMonitors;
