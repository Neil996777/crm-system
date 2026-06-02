export type DuplicateMatch = {
  type: 'lead' | 'account' | 'contact';
  matchStrength: 'High' | 'Medium' | 'Low';
  safeSummary: string;
  visible: boolean;
  rule: string;
};

export type DuplicateWarningResult = {
  result: 'PossibleDuplicate' | 'NoDuplicate';
  warningToken?: string;
  normalizedFields: string[];
  matches: DuplicateMatch[];
  rules: string[];
};
