import { FlexRow, MaterialIcon } from "components/globals";
import React from "react";
import {
  ProfileContainer,
  ProfileName,
  TopBarProductName,
  TopBarWrapper,
} from "./styles";

const TopBar: React.FunctionComponent = () => {
  return (
    <TopBarWrapper>
      <FlexRow>
        <TopBarProductName>Hatchet</TopBarProductName>
        <ProfileContainer>
          <ProfileName>Alexander Belanger</ProfileName>
          <MaterialIcon className="material-icons">
            arrow_drop_down
          </MaterialIcon>
        </ProfileContainer>
      </FlexRow>
    </TopBarWrapper>
  );
};

export default TopBar;
