import {
  SmallSpan,
  FlexRowLeft,
  altTextFontStack,
  SectionCard,
  FlexColScroll,
  FlexRowRight,
} from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

export const StatusText = styled(SmallSpan)`
  font-weight: bold;
`;

export const StatusContainer = styled(FlexRowLeft)<{ width?: string }>`
  ${altTextFontStack}
  color: ${(props) => props.theme.text.default};
  font-size: 13px;
  border: ${(props) => props.theme.line.default};
  padding: 2px 8px;
  border-radius: 5px;
  cursor: default;
  width: ${(props) => props.width || "fit-content"};
  :hover {
    background: ${(props) => props.theme.bg.hover};
  }

  > i {
    color: ${(props) => props.theme.text.default};
    font-size: 16px;
    margin-right: 10px;
  }
`;

export const TriggerContainer = styled(FlexRowLeft)`
  gap: 4px;
`;

export const StatusAndCommitContainer = styled(FlexRowRight)`
  gap: 8px;
`;

export const TriggerPRContainer = styled(StatusContainer)`
  font-size: 12px;
  padding: 4px 8px;
  cursor: pointer;

  > i {
    font-size: 12px;
  }
`;

export const GithubRefContainer = styled(StatusContainer)`
  padding: 6px 10px;
  cursor: pointer;
`;

export const GithubImg = styled.img`
  height: 18px;
  width: 18px;
  margin-right: 12px;
`;

export const ExpandedRunContainer = styled.div`
  overflow-x: auto;
  height: calc(100% - 200px);
`;

export const RunSectionCard = styled(SectionCard)`
  padding: 20px;
  height: fit-content;
`;
