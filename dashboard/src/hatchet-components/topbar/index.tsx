import { FlexRow, FlexRowLeft } from "../globals";
import React from "react";
import { TopBarProductName, TopBarWrapper } from "./styles";
import Logo from "../logo";

type Props = {
  is_authenticated?: boolean;
  children?: React.ReactNode;
};

const TopBar: React.FunctionComponent<Props> = ({
  is_authenticated = true,
  children,
}) => {
  return (
    <TopBarWrapper is_authenticated={is_authenticated}>
      <FlexRow>
        <FlexRowLeft>
          <Logo height="36px" width="36px" padding="6px" />
          <TopBarProductName>Hatchet</TopBarProductName>
        </FlexRowLeft>
        {children}
      </FlexRow>
    </TopBarWrapper>
  );
};

export default TopBar;
