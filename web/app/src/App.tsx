import React, { useState, useEffect, useRef } from "react";
import { createGlobalStyle, ThemeProvider } from "styled-components";
import { useSpring, interpolate } from "react-spring";
import { useGesture } from "react-with-gesture";

import { Drawer, DrawerItems, DrawerItem } from "@etxe/ui";
import { Header, HeaderItem} from "@etxe/ui";
import { Footer, FooterItem} from "@etxe/ui";
import { Overlay} from "@etxe/ui";

import { Main } from './components/Main'
import { Body } from './components/Body'
import { ActionHandler} from "./components/ActionHandler";
import { Cell} from "./components/Cell";
import {Content} from "./components/Content";

import MenuOpenIcon from '@material-ui/icons/MenuOpen'
import MenuIcon from '@material-ui/icons/Menu';

import AccountBoxIcon from "@material-ui/icons/AccountBox";
import DescriptionIcon from '@material-ui/icons/Description';
import LanguageIcon from '@material-ui/icons/Language';
import AddLocation from '@material-ui/icons/AddLocation';
import SignalWifi3Bar from '@material-ui/icons/SignalWifi3Bar';
import QuestionAnswerIcon from '@material-ui/icons/QuestionAnswer';
import SettingsIcon from '@material-ui/icons/Settings';

const drawerData = [
  { name: "User", url: "/private/loans", component: <AccountBoxIcon /> },
  { name: "Files upload", url: "/private/cards", component: <DescriptionIcon /> },
  { name: "Route", url: "/private/deposits", component: <LanguageIcon /> },
  { name: "Locations", url: "/private/services", component: <AddLocation /> },
  { name: "Connectivity", url: "/private/services", component: <SignalWifi3Bar /> },
  { name: "Social", url: "/private/services", component: <QuestionAnswerIcon /> },
  { name: "Settings", url: "/private/services", component: <SettingsIcon /> }
];


const GlobalStyle = createGlobalStyle`
  body {
    font-family: roboto;
    margin: 0px;
    color: #d8dcde;
  }
`;

export const App = ({ overlayColor = "transparent", width = 600 }) => {
  const [currentItem, setCurrentItem] = useState("User");

  const node = useRef(null);

  const handleClick = (event: Event) => {
    if (node.current.contains(event.target)) {
      return;
    }
    setActive(false);
  };

  // set the active state (true by default)
  const [active, setActive] = useState(true);

  // use react-with-gestures hook
  const [handler, { delta: [xDelta], down }] = useGesture();

  //gesture
  const { x, delta } = useSpring({
    native: true,
    to: {
      x: down ? xDelta : 0,
      delta: active ? 0 : -width
    },
    immediate: () => down
  });

  // drawer
  const { offset, color, backgroundColor } = useSpring({
    native: true,
    to: {
      offset: active ? 0 : 160,
      color: active ? "#344955" : "white",
      backgroundColor: active ? "#06151CA0" : "#FFFFFF00"
    }
  });

  useEffect(() => {
    document.addEventListener("mousedown", handleClick);
    return () => {
      document.removeEventListener("mousedown", handleClick);
    };
  }, []);

  useEffect(
    () => {
      if (!down && xDelta !== 0) {
        if (active && xDelta < -(width / 2)) {
          // when active, set the state back to inactive if dragged left for more than 1/2 of the width
          setActive(false);
        } else if (!active && xDelta > width / 4) {
          // when inactive, set the state back to active if dragged right for more than 1/4 of the width
          setActive(true);
        }
      }
    },
    [down] // trigger the effect of when down changes
  );

  return (
    <ThemeProvider theme={{ fontFamily: "roboto" }}>
      <React.Fragment>
        <Main>
          <Drawer
            ref={node}
            {...handler}
            width={width}
            style={{
              transform: interpolate(
                [x, delta],
                (x, delta) => `translateX(${Math.min(0, x + delta)}px)`
              )
            }}
          >
            <ActionHandler
              onClick={() => {
                if (xDelta !== 0) {
                  // prevent click if dragging
                  return;
                }
                setActive(!active);
              }}
              style={{
                color: color,
                transform: offset.interpolate(v => `translateX(${v}%)`)
              }}
            >
              {active ? <MenuOpenIcon /> : <MenuIcon />}
            </ActionHandler>
            <DrawerItems>
              {drawerData.map((item, index) => (
                <DrawerItem key={index}>
                  <Cell
                    current={currentItem === item.name}
                    onClick={() => { setCurrentItem(item.name); }}
                  >
                    {item.component}
                    {item.name}
                  </Cell>
                </DrawerItem>
              ))}
            </DrawerItems>
          </Drawer>
          <Overlay style={{ backgroundColor }} />
          <Content>
            <Header>
              <HeaderItem>{currentItem}</HeaderItem>
            </Header>
            <Body>
              <button
                onClick={() => {
                  if (xDelta !== 0) {
                    // prevent click if dragging
                    return;
                  }
                  setActive(!active);
                }}
                style={{
                  color: color,
                  transform: offset.interpolate(v => `translateX(${v}%)`)
                }}
              >
                click me
              </button>
            </Body>
            <Footer>
              <FooterItem>V1.0.0 alpha</FooterItem>
            </Footer>
          </Content>
        </Main>
        <GlobalStyle />
      </React.Fragment>
    </ThemeProvider>
  );
};
