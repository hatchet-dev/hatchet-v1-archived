import {
  MaterialIcon,
  FlexRow,
  StandardButton,
} from "hatchet-components";
import React from "react";
import {
  DeployedInfo,
  Version,
  VersionListWrapper,
  VersionWrapper,
} from "./styles";

export type Version = {
  version: string;
  link: string;
};

export type Props = {
  versions: Version[];
};

const DetailedVersionList: React.FC<Props> = ({ versions }) => {
  return (
    <VersionListWrapper>
      {versions.map((val, i) => {
        return (
          <VersionWrapper>
            <FlexRow>
              <a href={val.link} target="_blank">
                <Version>
                  <div>{val.version}</div>
                  <MaterialIcon className="material-icons">launch</MaterialIcon>
                </Version>
              </a>
              <DeployedInfo>Deployed on January 6th, 2022</DeployedInfo>
            </FlexRow>
            <StandardButton label="Deprecate" size="small" />
          </VersionWrapper>
        );
      })}
    </VersionListWrapper>
  );
};

export default DetailedVersionList;
