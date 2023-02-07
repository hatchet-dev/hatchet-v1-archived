import { MaterialIcon, Spinner } from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { ModuleRun } from "shared/api/generated/data-contracts";
import { relativeDate } from "shared/utils";
import ExpandedRun from "../../views/expandedmodule/components/expandedrun";
import {
  RunListWrapper,
  RunStatusWrapper,
  DateWrapper,
  RunWrapper,
  NextIconContainer,
} from "./styles";

export type Props = {
  runs: ModuleRun[];
};

const RunsList: React.FC<Props> = ({ runs }) => {
  const [selectedRun, setSelectedRun] = useState(null);

  if (selectedRun) {
    return (
      <RunListWrapper>
        <ExpandedRun back={() => setSelectedRun(null)} />
      </RunListWrapper>
    );
  }

  const renderIcon = (run: ModuleRun) => {
    switch (run.status) {
      case "completed":
        return <MaterialIcon className="material-icons">check</MaterialIcon>;
      case "failed":
        return <MaterialIcon className="material-icons">error</MaterialIcon>;
      case "in_progress":
        return <Spinner />;
    }
  };

  return (
    <RunListWrapper>
      {runs.map((val, i) => {
        return (
          <RunWrapper onClick={() => setSelectedRun(val)}>
            <RunStatusWrapper>
              {renderIcon(val)}
              <div>{val.status_description}</div>
            </RunStatusWrapper>
            <DateWrapper>
              <div>{relativeDate(val.created_at)}</div>
            </DateWrapper>
            <NextIconContainer>
              <MaterialIcon className="material-icons next-icon">
                navigate_next
              </MaterialIcon>
            </NextIconContainer>
          </RunWrapper>
        );
      })}
    </RunListWrapper>
  );
};

export default RunsList;
