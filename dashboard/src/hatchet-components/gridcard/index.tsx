import { fadeInAnimation } from "../globals";
import styled from "styled-components";

const GridCard = styled.div`
  ${fadeInAnimation}
  border: ${(props) => props.theme.line.default};
  align-items: center;
  user-select: none;
  border-radius: 8px;
  display: flex;
  font-size: 13px;
  font-weight: 500;
  padding: 3px 0px 5px;
  justify-content: center;
  cursor: pointer;
  color: ${(props) => props.theme.text.default};
  position: relative;
  background: ${(props) => props.theme.bg.shadeone};
  box-shadow: 0 4px 15px 0px #00000044;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  :hover {
    background: ${(props) => props.theme.bg.hover};
  }
`;

export default GridCard;
