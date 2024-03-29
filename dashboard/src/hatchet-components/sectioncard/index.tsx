import styled from "styled-components";

const SectionCard = styled.div`
  margin-bottom: 3px;
  border-radius: 6px;
  background: ${(props) => props.theme.bg.shadetwo};
  padding: 20px;
  overflow-y: auto;
  min-height: 180px;
`;

export default SectionCard;
