import {
  fadeInAnimation,
  FlexCol,
  FlexRow,
  textFontStack,
} from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const RunListWrapper = styled(FlexCol)`
  margin: 30px 0;
`;

export const RunStatusWrapper = styled(FlexRow)`
  width: 100%;
  font-size: 13px;
  justify-content: flex-start;

  > div {
    margin-left: 20px;
  }
`;

export const DateWrapper = styled(FlexRow)`
  min-width: 200px;
  font-size: 13px;
  color: ${(props) => props.theme.text.default};
`;

export const RunWrapper = styled(FlexRow)`
  ${fadeInAnimation}
  ${textFontStack}
  border: ${(props) => props.theme.line.default};
  background: ${(props) => props.theme.bg.shadeone};
  color: ${(props) => props.theme.text.default};
  cursor: pointer;
  margin-bottom: 3px;
  border-radius: 10px;
  padding: 14px;
  overflow: hidden;
  height: 60px;

  .next-icon {
    display: none;
    color: ${(props) => props.theme.text.default};
  }

  :hover .next-icon {
    display: inline-block;
  }
`;

export const NextIconContainer = styled.div`
  width: 30px;
  padding-top: 2px;

  > i {
    font-size: 18px;
  }
`;
