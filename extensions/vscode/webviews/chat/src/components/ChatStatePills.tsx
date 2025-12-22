import React from 'react';
import { Badge } from "@/components/ui/badge";
import { Bot, Drama, TerminalSquare, FileText, FileCodeCorner, FilePen } from 'lucide-react';
import { ToolTipper } from 'veg-webview-common';

type ChatStatePillsProps = {
  userInput: any;
  session: any;
};

export const ChatStatePills: React.FC<ChatStatePillsProps> = ({ userInput, session }) => {
  const cacheKeys = Object.keys(session?.state || {}).filter(k => k.startsWith("cache:"));
  const fileKeys = Object.keys(session?.state || {}).filter(k => k.startsWith("files:"));
  const agentmdKeys = Object.keys(session?.state || {}).filter(k => k.startsWith("agentmd:"));

  return (
    <>
      {/* User Settings */}
      {userInput?.agent && (
        <Badge className="text-sky-300 bg-sky-600/50">
          <Bot size={12} />
          {userInput?.agent}
        </Badge>
      )}
      {userInput?.model && (
        <Badge className="text-sky-300 bg-sky-600/50">
          <Drama size={12} />
          {userInput?.model}
        </Badge>
      )}
      {userInput?.environ && (
        <Badge className="text-lime-300/80 bg-lime-600/50">
          <TerminalSquare size={12} />
          {userInput?.environ}
        </Badge>
      )}

      {/* Context Info */}
      {agentmdKeys.length > 0 && (
        <ToolTipper label={agentmdKeys.map(k => k.split(':').slice(2).join(':')).join('\n')}>
          <Badge className="text-amber-200 bg-yellow-600/80 flex gap-1 items-center px-2">
            <FileText size={12} />
            <span>{agentmdKeys.length}</span>
          </Badge>
        </ToolTipper>
      )}
      {fileKeys.length > 0 && (
        <ToolTipper label={fileKeys.map(k => k.split(':').slice(2).join(':')).join('\n')}>
          <Badge className="text-violet-300 bg-violet-600/50 flex gap-1 items-center px-2">
            <FileCodeCorner size={12} />
            <span>{fileKeys.length}</span>
          </Badge>
        </ToolTipper>
      )}
      {cacheKeys.length > 0 && (
        <ToolTipper label={cacheKeys.map(k => k.split(':').slice(2).join(':')).join('\n')}>
          <Badge className="text-fuchsia-200/80 bg-fuchsia-600/50 flex gap-1 items-center px-2">
            <FilePen size={12} />
            <span>{cacheKeys.length}</span>
          </Badge>
        </ToolTipper>
      )}
    </>
  );
};
