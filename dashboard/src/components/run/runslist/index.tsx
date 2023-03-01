import { MaterialIcon, Spinner } from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { ModuleRun } from "shared/api/generated/data-contracts";
import { relativeDate } from "shared/utils";
import ExpandedRun from "../expandedrun";
import {
  RunListWrapper,
  RunStatusWrapper,
  DateWrapper,
  RunWrapper,
  NextIconContainer,
} from "./styles";

export type Props = {
  runs: ModuleRun[];
  select_run: (run: ModuleRun) => void;
};

const RunsList: React.FC<Props> = ({ runs, select_run }) => {
  const renderIcon = (run: ModuleRun) => {
    switch (run.status) {
      case "completed":
        return <MaterialIcon className="material-icons">check</MaterialIcon>;
      case "failed":
        return <MaterialIcon className="material-icons">error</MaterialIcon>;
      case "in_progress":
        return <Spinner />;
      case "queued":
        return <MaterialIcon className="material-icons">sort</MaterialIcon>;
    }
  };

  return (
    <RunListWrapper>
      {runs.map((val, i) => {
        return (
          <RunWrapper onClick={() => select_run(val)}>
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
