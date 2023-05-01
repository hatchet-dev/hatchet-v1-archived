import { FlexCol, FlexRow, FlexRowLeft } from "hatchet-components";
import styled from "styled-components";

export const VersionListWrapper = styled(FlexCol)`
  margin: 10px 0;
`;

export const VersionWrapper = styled(FlexRow)`
  border: ${(props) => props.theme.line.default};
  border-radius: 6px;
  margin: 6px 0;
  height: 60px;
`;

export const VersionMetadataContainer = styled(FlexRow)``;

export const Version = styled(FlexRowLeft)<{ selected?: boolean }>`
  color: ${(props) => props.theme.text.default};
  padding: 2px 4px 0 4px;
  margin: 0 10px;
  min-width: 80px;
  cursor: pointer;

  > div {
    font-size: 14px;
    font-weight: 500;
  }

  > i {
    display: none;
    margin-left: 8px;
  }

  :hover {
    border-bottom: ${(props) => props.theme.line.selected};
    margin-bottom: -1px;
  }

  :hover > i {
    display: inline-block;
  }
`;

export const DeployedInfo = styled.div`
  font-size: 13px;
  color: ${(props) => props.theme.text.default};
`;
