"use client";

import { useState } from "react";
import Image from "next/image";
import { MoreVertical, Trash } from "lucide-react";
import { Button } from "@/_components/atoms/button";
import {
  Card as CardUI,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/_components/atoms/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/_components/atoms/dropdown-menu";

interface Collection {
  id: string;
  name: string;
  description: string;
  numberOfCommands: number;
  imageUrl: string;
}

const Card = ({
  collection,
  onDelete,
}: {
  collection: Collection;
  onDelete: (id: string) => void;
}) => {
  return (
    <CardUI className="aspect-video rounded-xl bg-muted/50">
      <CardHeader className="relative">
        <Image
          src={collection.imageUrl}
          alt={collection.name}
          width={300}
          height={200}
          className="h-48 w-full rounded-t-lg object-cover"
        />
        <div className="absolute right-2 top-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => onDelete(collection.id)}>
                <Trash className="mr-2 h-4 w-4" />
                <span>Delete</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>
      <CardContent>
        <CardTitle className="mb-1">{collection.name}</CardTitle>
        <p className="mb-2 text-sm text-gray-600">{collection.description}</p>
      </CardContent>
      <CardFooter>
        <p className="text-sm text-gray-500">
          Commands: {collection.numberOfCommands}
        </p>
      </CardFooter>
    </CardUI>
  );
};

export function CardList() {
  const [collections, setCollections] = useState<Collection[]>([
    {
      id: "1",
      name: "Web Development",
      description: "A collection of web development commands",
      numberOfCommands: 25,
      imageUrl: "/placeholder.svg?height=200&width=300",
    },
    {
      id: "2",
      name: "Data Science",
      description:
        "Essential commands for data analysisEssential commands for data analysisEssential commands for data analysis",
      numberOfCommands: 30,
      imageUrl: "/placeholder.svg?height=200&width=300",
    },
    {
      id: "3",
      name: "DevOps",
      description: "Commands for streamlining DevOps processes",
      numberOfCommands: 20,
      imageUrl: "/placeholder.svg?height=200&width=300",
    },
    {
      id: "4",
      name: "DevOps",
      description: "Commands for streamlining DevOps processes",
      numberOfCommands: 20,
      imageUrl: "/placeholder.svg?height=200&width=300",
    },
    {
      id: "5",
      name: "DevOps",
      description: "Commands for streamlining DevOps processes",
      numberOfCommands: 20,
      imageUrl: "/placeholder.svg?height=200&width=300",
    },
  ]);

  const handleDelete = (id: string) => {
    setCollections((prevCollections) =>
      prevCollections.filter((collection) => collection.id !== id),
    );
  };

  return (
    <>
      <h1 className="mb-6 text-2xl font-bold">Collections</h1>
      <div className="grid auto-rows-min gap-4 md:grid-cols-4 lg:grid-cols-5">
        {/* <div className="aspect-video rounded-xl bg-muted/50" /> */}
        {/* <div className="aspect-video rounded-xl bg-muted/50" /> */}
        {/* <div className="aspect-video rounded-xl bg-muted/50" /> */}
        {/* <div className="aspect-video rounded-xl bg-muted/50" /> */}

        {collections.map((collection) => (
          <Card
            key={collection.id}
            collection={collection}
            onDelete={handleDelete}
          />
        ))}
      </div>
    </>
  );
}
