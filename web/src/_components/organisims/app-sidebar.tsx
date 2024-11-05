"use client";

import * as React from "react";
import { Folders, Tags, Terminal } from "lucide-react";

import { NavMain } from "@/_components/organisims/nav-main";
import { NavUser } from "@/_components/organisims/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarRail,
} from "@/_components/atoms/sidebar";

const data = {
  user: {
    name: "User",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      name: "Commands",
      url: "#",
      icon: Terminal,
    },
    {
      name: "Collections",
      url: "/dashboard/collections",
      icon: Folders,
    },
    {
      name: "Tags",
      url: "#",
      icon: Tags,
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      {/* <SidebarHeader> */}
      {/*   <TeamSwitcher teams={data.teams} /> */}
      {/* </SidebarHeader> */}
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
