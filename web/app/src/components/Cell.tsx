import styled from "styled-components";
import { lighten } from "polished";

import { resolveBackgroundColor, resolveColor } from "../styles/helpers";

export const Cell = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  ${resolveColor};
  ${resolveBackgroundColor}

  padding-top: 10px;
  padding-bottom: 10px;
  padding-left: 10px;
  :hover {
    background-color: ${() => lighten(0.5, "#5f7481")};
  }
`;
