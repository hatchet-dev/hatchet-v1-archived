import github from "assets/github.png";
import React, { useState } from "react";
import { StatusText } from "components/status/styles";
import {
  FlexRow,
  FlexCol,
  MaterialIcon,
} from "@hatchet-dev/hatchet-components";
import {
  ExpandableSettingsContainer,
  ExpandableSettingsHeader,
  ExpandableSettingsText,
  ExpandableSettingsBody,
} from "./styles";

type Props = {
  text: string;
  children?: React.ReactNode;
};

const ExpandableSettings: React.FC<Props> = ({ text, children }) => {
  const [expanded, setExpanded] = useState(false);
  return (
    <ExpandableSettingsContainer>
      <ExpandableSettingsHeader onClick={() => setExpanded(!expanded)}>
        <ExpandableSettingsText>{text}</ExpandableSettingsText>
        <MaterialIcon className="material-icons">
          {expanded ? "expand_less" : "expand_more"}
        </MaterialIcon>
      </ExpandableSettingsHeader>
      {expanded && <ExpandableSettingsBody>{children}</ExpandableSettingsBody>}
    </ExpandableSettingsContainer>
  );
};

export default ExpandableSettings;
