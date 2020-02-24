# ROW FIELD VALIDATION

## Fields list

- ParentDirectory
- Filename
- SizeHR
- SizeInBytes
- Checksum
- RemoveFile

## ParentDirectory

- Required
- Valid path
- Type conversion: NO

## Filename

- Required
- Type conversion: NO

## SizeHR

- Optional?
  - User could have chosen to remove it and leave the field empty
- Type conversion: NO

## SizeInBytes

- Optional?
  - User could have chosen to remove it and leave the field empty
  - We can recalculate this later if missing
- Type conversion: Yes
  - we will likely wish to provide a summary of space reclaimed at the end

## Checksum

- Required
- Used to validate that we are going to remove the correct file
- Type conversion: NO

## RemoveFile

- Optional
  - User might not wish to remove a file
- Type conversion: YES
