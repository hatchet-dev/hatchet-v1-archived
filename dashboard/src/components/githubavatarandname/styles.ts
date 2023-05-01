import { SmallSpan, FlexRowLeft } from "hatchet-components";
import styled from "styled-components";

export const GithubAvatar = styled.img<{ avatar_size: "default" | "small" }>`
  height: ${(props) => (props.avatar_size == "default" ? "24px" : "18px")};
  width: ${(props) => (props.avatar_size == "default" ? "24px" : "18px")};
  border-radius: 100%;
`;

export const GithubName = styled(SmallSpan)`
  font-weight: bold;
`;

export const Container = styled(FlexRowLeft)`
  gap: 8px;
`;
