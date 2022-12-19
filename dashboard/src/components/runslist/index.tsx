import { MaterialIcon } from "components/globals";
import React, { useState } from "react";
import ExpandedRun from "../../views/expandedmodule/components/expandedrun";
import {
  RunListWrapper,
  RunStatusWrapper,
  DateWrapper,
  RunWrapper,
  NextIconContainer,
} from "./styles";

export type Run = {
  status: string;
  date: string;
};

export type Props = {
  runs: Run[];
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

  return (
    <RunListWrapper>
      {runs.map((val, i) => {
        return (
          <RunWrapper onClick={() => setSelectedRun(val)}>
            <RunStatusWrapper>
              <MaterialIcon className="material-icons">error</MaterialIcon>
              <div>Successfully deployed infrastructure.</div>
            </RunStatusWrapper>
            <DateWrapper>
              <div>{val.date}</div>
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
